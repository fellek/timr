<?php

namespace TimetrackerBundle\Controller;

use Symfony\Component\HttpFoundation\Request;
use Symfony\Bundle\FrameworkBundle\Controller\Controller;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Method;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Route;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\Template;
use TimetrackerBundle\Entity\Employee;
use TimetrackerBundle\Form\EmployeeType;
use Symfony\Component\Security\Core\SecurityContext;
use Symfony\Component\Security\Core\Exception\AccessDeniedException;

/**
 * Employee controller.
 *
 * @Route("/employee")
 */
class EmployeeController extends Controller
{

    /**
     * Lists all Employee entities.
     *
     * @Route("/", name="employee")
     * @Method("GET")
     * @Template()
     */
    public function indexAction()
    {
        $em = $this->getDoctrine()->getManager();

		$sc = $this->get('security.context');
		$user = $sc->getToken()->getUser();

		$employees = array();
		if( $sc->isGranted('ROLE_ADMIN') ) {
        	$employees = $em->getRepository('TimetrackerBundle:Employee')->findAll();
		} else if( $sc->isGranted('ROLE_MDBOSS') ) {
			$allEmployees = $em->getRepository('TimetrackerBundle:Employee')->findAll();
			foreach( $allEmployees as $employee ) {
				$allRoles = $employee->getRoles();
				foreach( $allRoles as $role ) {
					if( $role->getRole() == 'ROLE_MEDIADESIGNER' ) {
						$employees[] = $employee;
					}
				}
			}
		} else {
			$employees[] = $em->getRepository('TimetrackerBundle:Employee')->find($user->getId());
		}

        return array(
            'employees' => $employees,
        );
    }
    /**
     * Creates a new Employee entity.
     *
     * @Route("/", name="employee_create")
     * @Method("POST")
     * @Template("TimetrackerTimetrackerBundle:Employee:new.html.twig")
     */
    public function createAction(Request $request)
    {
        $entity = new Employee();
        $form = $this->createCreateForm($entity);
        $form->handleRequest($request);

        if ($form->isValid()) {
     		$encoder = $this->get('security.encoder_factory')->getEncoder($entity);
			$hashed_password = $encoder->encodePassword($entity->getPassword(), $entity->getSalt());
    		$entity->setPassword($hashed_password);
    		
            $em = $this->getDoctrine()->getManager();
			$role = $em->getRepository('TimetrackerBundle:Role')->findOneByRole('ROLE_USER');
			$entity->addRole($role);

            $em->persist($entity);
            $em->flush();

            return $this->redirect($this->generateUrl('employee_show', array('id' => $entity->getId())));
        }

        return array(
            'entity' => $entity,
            'form'   => $form->createView(),
        );
    }

    /**
     * Creates a form to create a Employee entity.
     *
     * @param Employee $entity The entity
     *
     * @return \Symfony\Component\Form\Form The form
     */
    private function createCreateForm(Employee $entity)
    {
        $form = $this->createForm(new EmployeeType(), $entity, array(
            'action' => $this->generateUrl('employee_create'),
            'method' => 'POST',
        ));

        return $form;
    }

    /**
     * Displays a form to create a new Employee entity.
     *
     * @Route("/new", name="employee_new")
     * @Method("GET")
     * @Template()
     */
    public function newAction()
    {
        $entity = new Employee();
        $form   = $this->createCreateForm($entity);

        return array(
            'entity' => $entity,
            'form'   => $form->createView(),
        );
    }

    /**
     * Finds and displays a Employee entity.
     *
     * @Route("/{id}", name="employee_show")
     * @Method("GET")
     * @Template()
     */
    public function showAction($id)
    {
        $em = $this->getDoctrine()->getManager();

        $employee = $em->getRepository('TimetrackerBundle:Employee')->find($id);

        if (!$employee) {
            throw $this->createNotFoundException('Unable to find Employee entity.');
        }

        $deleteForm = $this->createDeleteForm($id);

        return array(
            'employee'    => $employee,
            'delete_form' => $deleteForm->createView(),
        );
    }

    /**
     * Displays a form to edit an existing Employee entity.
     *
     * @Route("/{id}/edit", name="employee_edit")
     * @Method("GET")
     * @Template()
     */
    public function editAction($id)
    {
    	$sc = $this->get('security.context');
		$user = $sc->getToken()->getUser();
		$userId = $user->getId();

        if( $id != $userId && !$sc->isGranted('ROLE_ADMIN') && !$sc->isGranted('ROLE_MDBOSS') ) {
        	throw new AccessDeniedException("Unauthorized Access");
        } else if( $sc->isGranted('ROLE_MDBOSS') && $id != $userId ) {
        	$allRoles = $employee->getRoles();
        	$show = 0;
        	foreach( $allRoles as $role ) {
        		if( $role->getRole() == 'ROLE_MEDIADESIGNER' ) {
        			$show = 1;
        		}
        	}
        	if( $show == 0 ) {
        		throw new AccessDeniedException("Unauthorized Access");
        	}
        }

        $em = $this->getDoctrine()->getManager();

        $employee = $em->getRepository('TimetrackerBundle:Employee')->find($id);

        if (!$employee) {
            throw $this->createNotFoundException('Unable to find Employee entity.');
        }

        $editForm = $this->createEditForm($employee);
        $deleteForm = $this->createDeleteForm($id);

        return array(
            'employee'    => $employee,
            'edit_form'   => $editForm->createView(),
            'delete_form' => $deleteForm->createView(),
        );
    }

    /**
    * Creates a form to edit a Employee entity.
    *
    * @param Employee $entity The entity
    *
    * @return \Symfony\Component\Form\Form The form
    */
    private function createEditForm(Employee $entity)
    {
        $form = $this->createForm(new EmployeeType(), $entity, array(
            'action' => $this->generateUrl('employee_update', array('id' => $entity->getId())),
            'method' => 'PUT',
        ));

        return $form;
    }
    /**
     * Edits an existing Employee entity.
     *
     * @Route("/{id}", name="employee_update")
     * @Method("PUT")
     * @Template("TimetrackerTimetrackerBundle:Employee:edit.html.twig")
     */
    public function updateAction(Request $request, $id)
    {
        $em = $this->getDoctrine()->getManager();

        $entity = $em->getRepository('TimetrackerBundle:Employee')->find($id);

        if (!$entity) {
            throw $this->createNotFoundException('Unable to find Employee entity.');
        }

        $deleteForm = $this->createDeleteForm($id);
        $editForm = $this->createEditForm($entity);
        $editForm->handleRequest($request);

        if ($editForm->isValid()) {
     		$encoder = $this->get('security.encoder_factory')->getEncoder($entity);
			$hashed_password = $encoder->encodePassword($entity->getPassword(), $entity->getSalt());
    		$entity->setPassword($hashed_password);
 
            $em->flush();

            return $this->redirect($this->generateUrl('employee_edit', array('id' => $id)));
        }

        return array(
            'entity'      => $entity,
            'edit_form'   => $editForm->createView(),
            'delete_form' => $deleteForm->createView(),
        );
    }
    /**
     * Deletes a Employee entity.
     *
     * @Route("/{id}", name="employee_delete")
     * @Method("DELETE")
     */
    public function deleteAction(Request $request, $id)
    {
        $form = $this->createDeleteForm($id);
        $form->handleRequest($request);

        if ($form->isValid()) {
            $em = $this->getDoctrine()->getManager();
            $entity = $em->getRepository('TimetrackerBundle:Employee')->find($id);

            if (!$entity) {
                throw $this->createNotFoundException('Unable to find Employee entity.');
            }

            $em->remove($entity);
            $em->flush();
        }

        return $this->redirect($this->generateUrl('employee'));
    }

    /**
     * Creates a form to delete a Employee entity by id.
     *
     * @param mixed $id The entity id
     *
     * @return \Symfony\Component\Form\Form The form
     */
    private function createDeleteForm($id)
    {
        return $this->createFormBuilder()
            ->setAction($this->generateUrl('employee_delete', array('id' => $id)))
            ->setMethod('DELETE')
            ->getForm();
    }
}